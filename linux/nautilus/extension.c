#include <glib.h>

#include <libnautilus-extension/nautilus-extension-types.h>
#include <libnautilus-extension/nautilus-file-info.h>
#include <libnautilus-extension/nautilus-info-provider.h>
#include <libnautilus-extension/nautilus-menu-provider.h>
#include <libnautilus-extension/nautilus-property-page-provider.h>

#include <gtk/gtktable.h>
#include <gtk/gtkvbox.h>
#include <gtk/gtkhbox.h>
#include <gtk/gtklabel.h>
#include <string.h>
#include <time.h>

static GType provider_types[1];
static GType share_link_extension_type;
static GObjectClass *parent_class;

typedef struct {
	GObject parent_slot;
} ShareLinkExtension;

typedef struct {
	GObjectClass parent_slot;
} ShareLinkExtensionClass;

void nautilus_module_initialize (GTypeModule  *module);
void nautilus_module_shutdown (void);
void nautilus_module_list_types (const GType **types, int *num_types);
GType share_link_extension_get_type (void);

static void share_link_extension_register_type (GTypeModule *module);

static GList * share_link_extension_get_file_items (NautilusMenuProvider *provider,
                GtkWidget *window,
                GList *files);
#if 0
static GList * share_link_extension_get_background_items (NautilusMenuProvider *provider,
                GtkWidget *window,
                NautilusFileInfo *current_folder);
static GList * share_link_extension_get_toolbar_items (NautilusMenuProvider *provider,
                GtkWidget *window,
                NautilusFileInfo *current_folder);
#endif

static void share_link (NautilusMenuItem *item, gpointer user_data);

void nautilus_module_initialize (GTypeModule  *module)
{
        share_link_extension_register_type (module);

        provider_types[0] = share_link_extension_get_type ();
}

void nautilus_module_shutdown (void)
{

}

void nautilus_module_list_types (const GType **types,
                                 int *num_types)
{
        *types = provider_types;
        *num_types = G_N_ELEMENTS (provider_types);
}

GType share_link_extension_get_type (void)
{
        return share_link_extension_type;
}

static void share_link_extension_instance_init (ShareLinkExtension *object)
{
}

static void share_link_extension_class_init(ShareLinkExtensionClass *class)
{
	parent_class = g_type_class_peek_parent (class);
}

static void share_link_extension_menu_provider_iface_init(
		NautilusMenuProviderIface *iface)
{
	iface->get_file_items = share_link_extension_get_file_items;
}

static void share_link_extension_register_type (GTypeModule *module)
{
        static const GTypeInfo info = {
                sizeof (ShareLinkExtensionClass),
                (GBaseInitFunc) NULL,
                (GBaseFinalizeFunc) NULL,
                (GClassInitFunc) share_link_extension_class_init,
                NULL,
                NULL,
                sizeof (ShareLinkExtension),
                0,
                (GInstanceInitFunc) share_link_extension_instance_init,
        };

	static const GInterfaceInfo menu_provider_iface_info = {
		(GInterfaceInitFunc) share_link_extension_menu_provider_iface_init,
		NULL,
		NULL
	};

        share_link_extension_type = g_type_module_register_type (module,
                             G_TYPE_OBJECT,
                             "ShareLinkExtension",
                             &info, 0);

	g_type_module_add_interface (module,
				     share_link_extension_type,
				     NAUTILUS_TYPE_MENU_PROVIDER,
				     &menu_provider_iface_info);
}


static GList * share_link_extension_get_file_items (NautilusMenuProvider *provider,
                GtkWidget *window,
                GList *files)
{
        NautilusMenuItem *item;
        GList *l;
        GList *ret;

#if 0
        /* This extension only operates on selections that include only
         * share_link files */
        for (l = files; l != NULL; l = l->next) {
                NautilusFileInfo *file = NAUTILUS_FILE_INFO (l->data);
                if (!nautilus_file_is_mime_type (file, "application/x-foo")) {
                        return;
                }
        }
#endif

        for (l = files; l != NULL; l = l->next) {
                NautilusFileInfo *file = NAUTILUS_FILE_INFO (l->data);
                char *name;
                name = nautilus_file_info_get_name (file);
                g_print ("selected %s\n", name);
                g_free (name);
        }

        item = nautilus_menu_item_new ("ShareLinkExtension::share_link_cb",
                                       "Share link via Storj",
                                       "Generates public link via Storj and copies ot to the clipboard",
                                       NULL /* icon name */);
        g_signal_connect (item, "activate", G_CALLBACK (do_stuff_cb), provider);
        g_object_set_data_full ((GObject*) item, "share_link_extension_files",
                                nautilus_file_info_list_copy (files),
                                (GDestroyNotify)nautilus_file_info_list_free);
        ret = g_list_append (NULL, item);

        return ret;
}

#if 0
static GList * share_link_extension_get_background_items (NautilusMenuProvider *provider,
                GtkWidget *window,
                NautilusFileInfo *current_folder)
{
        return NULL;
}

static GList * share_link_extension_get_toolbar_items (NautilusMenuProvider *provider,
                GtkWidget *window,
                NautilusFileInfo *current_folder)
{
        return NULL;
}
#endif

static void share_link_cb (NautilusMenuItem *item,
                         gpointer user_data)
{
        GList *files;
        GList *l;

        files = g_object_get_data ((GObject *) item, "share_link_extension_files");

        for (l = files; l != NULL; l = l->next) {
                NautilusFileInfo *file = NAUTILUS_FILE_INFO (l->data);
                char *name;
                name = nautilus_file_info_get_name (file);
                g_print ("generating link via storj %s\n", name);
                g_free (name);
        }
}