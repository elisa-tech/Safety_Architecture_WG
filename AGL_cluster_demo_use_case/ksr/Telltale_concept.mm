<map version="freeplane 1.8.0">
<!--To view this file, download free mind mapping software Freeplane from http://freeplane.sourceforge.net -->
<attribute_registry SHOW_ATTRIBUTES="hide"/>
<node TEXT="Safety concept Cluster Display use case" LOCALIZED_STYLE_REF="AutomaticLayout.level.root" FOLDED="false" ID="ID_273763478" CREATED="1609081280555" MODIFIED="1610894422559"><hook NAME="MapStyle" zoom="1.003">
    <properties show_icon_for_attributes="false" edgeColorConfiguration="#808080ff,#ff0000ff,#0000ffff,#00ff00ff,#ff00ffff,#00ffffff,#7c0000ff,#00007cff,#007c00ff,#7c007cff,#007c7cff,#7c7c00ff" show_note_icons="true" fit_to_viewport="false"/>

<map_styles>
<stylenode LOCALIZED_TEXT="styles.root_node" STYLE="oval" UNIFORM_SHAPE="true" VGAP_QUANTITY="24.0 pt">
<font SIZE="24"/>
<stylenode LOCALIZED_TEXT="styles.predefined" POSITION="right" STYLE="bubble">
<stylenode LOCALIZED_TEXT="default" ICON_SIZE="12.0 pt" COLOR="#000000" STYLE="fork">
<font NAME="SansSerif" SIZE="10" BOLD="false" ITALIC="false"/>
</stylenode>
<stylenode LOCALIZED_TEXT="defaultstyle.details"/>
<stylenode LOCALIZED_TEXT="defaultstyle.attributes">
<font SIZE="9"/>
</stylenode>
<stylenode LOCALIZED_TEXT="defaultstyle.note" COLOR="#000000" BACKGROUND_COLOR="#ffffff" TEXT_ALIGN="LEFT"/>
<stylenode LOCALIZED_TEXT="defaultstyle.floating">
<edge STYLE="hide_edge"/>
<cloud COLOR="#f0f0f0" SHAPE="ROUND_RECT"/>
</stylenode>
</stylenode>
<stylenode LOCALIZED_TEXT="styles.user-defined" POSITION="right" STYLE="bubble">
<stylenode LOCALIZED_TEXT="styles.topic" COLOR="#18898b" STYLE="fork">
<font NAME="Liberation Sans" SIZE="10" BOLD="true"/>
</stylenode>
<stylenode LOCALIZED_TEXT="styles.subtopic" COLOR="#cc3300" STYLE="fork">
<font NAME="Liberation Sans" SIZE="10" BOLD="true"/>
</stylenode>
<stylenode LOCALIZED_TEXT="styles.subsubtopic" COLOR="#669900">
<font NAME="Liberation Sans" SIZE="10" BOLD="true"/>
</stylenode>
<stylenode LOCALIZED_TEXT="styles.important">
<icon BUILTIN="yes"/>
</stylenode>
<stylenode TEXT="Teststyle" COLOR="#cc00cc" STYLE="wide_hexagon" NUMBERED="false">
<edge COLOR="#00ff00"/>
<cloud COLOR="#ffcc66" SHAPE="ARC"/>
<hook NAME="NodeConditionalStyles">
    <conditional_style ACTIVE="true" LOCALIZED_STYLE_REF="AutomaticLayout.level,1" LAST="false">
        <attribute_compare_condition VALUE="FSR" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
</hook>
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="FSR"/>
</stylenode>
<stylenode TEXT="Requirement">
<edge COLOR="#007c00"/>
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE=""/>
<hook NAME="NodeConditionalStyles">
    <conditional_style ACTIVE="true" STYLE_REF="ASIL A" LAST="false">
        <attribute_compare_condition VALUE="A" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL A[A]" LAST="false">
        <attribute_compare_condition VALUE="A[A]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL A[B]" LAST="false">
        <attribute_compare_condition VALUE="A[B]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL A[C]" LAST="false">
        <attribute_compare_condition VALUE="A[C]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL A[D]" LAST="false">
        <attribute_compare_condition VALUE="A[D]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL B" LAST="false">
        <attribute_compare_condition VALUE="B" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL B[B]" LAST="false">
        <attribute_compare_condition VALUE="B[B]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL B[C]" LAST="false">
        <attribute_compare_condition VALUE="B[C]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL B[D]" LAST="false">
        <attribute_compare_condition VALUE="B[D]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL C" LAST="false">
        <attribute_compare_condition VALUE="C" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL C[C]" LAST="false">
        <attribute_compare_condition VALUE="C[C]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL C[D]" LAST="false">
        <attribute_compare_condition VALUE="C[D]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL D" LAST="false">
        <attribute_compare_condition VALUE="D" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="ASIL D[D]" LAST="false">
        <attribute_compare_condition VALUE="D[D]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="QM" LAST="false">
        <attribute_compare_condition VALUE="QM" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="QM[A]" LAST="false">
        <attribute_compare_condition VALUE="QM[A]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="QM[B]" LAST="false">
        <attribute_compare_condition VALUE="QM[B]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="QM[C]" LAST="false">
        <attribute_compare_condition VALUE="QM[C]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="QM[D]" LAST="false">
        <attribute_compare_condition VALUE="QM[D]" ATTRIBUTE="ASIL" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="FSR" LAST="false">
        <attribute_compare_condition VALUE="FSR" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="TSR" LAST="false">
        <attribute_compare_condition VALUE="TSR" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="Info" LAST="false">
        <attribute_compare_condition VALUE="Information" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="SZ" LAST="false">
        <attribute_compare_condition VALUE="SZ" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="HW" LAST="false">
        <attribute_compare_condition VALUE="HW" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
    <conditional_style ACTIVE="true" STYLE_REF="SW" LAST="false">
        <attribute_compare_condition VALUE="SW" ATTRIBUTE="Type" COMPARATION_RESULT="0" SUCCEED="true"/>
    </conditional_style>
</hook>
</stylenode>
<stylenode TEXT="QM">
<icon BUILTIN="ASIL_QM"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="QM[A]">
<icon BUILTIN="ASIL_QM[A]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="QM[B]">
<icon BUILTIN="ASIL_QM[B]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="QM[C]">
<icon BUILTIN="ASIL_QM[C]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="QM[D]">
<icon BUILTIN="ASIL_QM[D]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL A">
<icon BUILTIN="ASIL_A"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL A[A]">
<icon BUILTIN="ASIL_A[A]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL A[B]">
<icon BUILTIN="ASIL_A[B]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL A[C]">
<icon BUILTIN="ASIL_A[C]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL A[D]">
<icon BUILTIN="ASIL_A[D]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL B">
<icon BUILTIN="ASIL_B"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL B[B]">
<icon BUILTIN="ASIL_B[B]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL B[C]">
<icon BUILTIN="ASIL_B[C]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL B[D]">
<icon BUILTIN="ASIL_B[D]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL C">
<icon BUILTIN="ASIL_C"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL C[C]">
<icon BUILTIN="ASIL_C[C]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL C[D]">
<icon BUILTIN="ASIL_C[D]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL D">
<icon BUILTIN="ASIL_D"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="ASIL D[D]">
<icon BUILTIN="ASIL_D[D]"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="SZ" STYLE="wide_hexagon">
<icon BUILTIN="SZ"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="FSR" STYLE="oval">
<icon BUILTIN="FSR"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="TSR" STYLE="bubble">
<icon BUILTIN="TSR"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="Info">
<icon BUILTIN="Info"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="HW" STYLE="rectangle" BORDER_COLOR_LIKE_EDGE="false" BORDER_COLOR="#ff9900">
<icon BUILTIN="HW"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="SW" STYLE="rectangle" BORDER_COLOR_LIKE_EDGE="false" BORDER_COLOR="#0066ff">
<icon BUILTIN="SW"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="Warning" BACKGROUND_COLOR="#ff0033" STYLE="rectangle" BORDER_COLOR_LIKE_EDGE="false" BORDER_COLOR="#0066ff">
<icon BUILTIN="button_cancel"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="Source code Tag" ICON_SIZE="12.0 pt" BACKGROUND_COLOR="#cccccc" STYLE="narrow_hexagon" BORDER_COLOR_LIKE_EDGE="false" BORDER_COLOR="#0066ff">
<icon BUILTIN="very_positive"/>
<font NAME="L M Mono Caps10" BOLD="true"/>
<edge COLOR="#007c00"/>
</stylenode>
<stylenode TEXT="KSR" STYLE="rectangle" BORDER_COLOR_LIKE_EDGE="false" BORDER_COLOR="#ff3333">
<hook NAME="NodeConditionalStyles">
    <conditional_style ACTIVE="true" STYLE_REF="Requirement" LAST="false">
        <attribute_contains_condition ATTRIBUTE="ASIL" VALUE="B"/>
    </conditional_style>
</hook>
<edge COLOR="#007c00"/>
</stylenode>
</stylenode>
<stylenode LOCALIZED_TEXT="styles.AutomaticLayout" POSITION="right" STYLE="bubble">
<stylenode LOCALIZED_TEXT="AutomaticLayout.level.root" COLOR="#000000" STYLE="oval" SHAPE_HORIZONTAL_MARGIN="10.0 pt" SHAPE_VERTICAL_MARGIN="10.0 pt">
<font SIZE="18"/>
</stylenode>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,1" COLOR="#0033ff">
<font SIZE="16"/>
</stylenode>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,2" COLOR="#00b439">
<font SIZE="14"/>
</stylenode>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,3" COLOR="#990000">
<font SIZE="12"/>
</stylenode>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,4" COLOR="#111111">
<font SIZE="10"/>
</stylenode>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,5"/>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,6"/>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,7"/>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,8"/>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,9"/>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,10"/>
<stylenode LOCALIZED_TEXT="AutomaticLayout.level,11"/>
</stylenode>
</stylenode>
</map_styles>
</hook>
<hook NAME="AutomaticEdgeColor" COUNTER="54" RULE="ON_BRANCH_CREATION"/>
<node TEXT="while requested, the system shall display the driver warning within 200 ms or transition to the safe state within 200 ms." STYLE_REF="Requirement" POSITION="right" ID="ID_971613141" CREATED="1609106418278" MODIFIED="1610897196058" HGAP_QUANTITY="19.99999982118607 pt" VSHIFT_QUANTITY="2.9999999105930186 pt">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SZ"/>
<attribute NAME="Allocation" VALUE=""/>
<node TEXT="Information: “while ” means that, if the telltale request persists/is repeated, the system has to continue to display the telltale." STYLE_REF="Requirement" ID="ID_1502911625" CREATED="1609087027838" MODIFIED="1610897196059">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="The 200 ms include the time needed for the request to reach the Cluster demo. This is considered in the frequency of the cyclic communication." STYLE_REF="Requirement" ID="ID_1052985289" CREATED="1609428730544" MODIFIED="1610897196061">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="The Telltale requester shall send a request cyclically controlling whether a telltale is needed to be shown or not." STYLE_REF="Requirement" ID="ID_1780168904" CREATED="1609154124607" MODIFIED="1610897196063">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="FSR"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
<node TEXT="The Telltale requester shall send the telltale request message every 200 ms" STYLE_REF="Requirement" ID="ID_736988533" CREATED="1610617420206" MODIFIED="1610897196066">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
</node>
<node TEXT="The Telltale request message shall contain a boolean &quot;telltale_request&quot; = 0 if the telltale is not requested and 1 if the telltale is requested" STYLE_REF="Requirement" ID="ID_529767340" CREATED="1610617528540" MODIFIED="1610897196068">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
</node>
<node TEXT="The Telltale request message shall be E2E protected with E2E Protocol xxx" STYLE_REF="Requirement" ID="ID_950923064" CREATED="1610617442339" MODIFIED="1610897196069">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
<node TEXT="We don&apos;t specify this in all detail here, Message counter and CRC is needed" STYLE_REF="Requirement" ID="ID_1340201467" CREATED="1610617507393" MODIFIED="1610897196070">
<attribute NAME="Type" VALUE="Information"/>
</node>
</node>
</node>
<node TEXT="All inputs from outside the system, the cluster controller uses to determine whether a requested telltale is shown shall be E2E protected against data corruption out of order transmission and message loss" STYLE_REF="Requirement" ID="ID_883554261" CREATED="1609428983812" MODIFIED="1610897196072">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="FSR"/>
<node TEXT="The Cluster controller shall monitor messages from the Telltale requester" STYLE_REF="Requirement" ID="ID_1807969240" CREATED="1610617925379" MODIFIED="1610897196075">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
</node>
<node TEXT="The Cluster controller shall check the telltale request message for E2E miss" STYLE_REF="Requirement" ID="ID_199781775" CREATED="1610617974658" MODIFIED="1610897196078">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
<node ID="ID_1822634618" TREE_ID="ID_1340201467"/>
</node>
<node TEXT="If the cluster controller determines an E2E miss in the tell tale request message, the cluster controller shall transition the system into the safe state" STYLE_REF="Requirement" ID="ID_1213070481" CREATED="1610618005875" MODIFIED="1610897196081">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
</node>
<node TEXT="The Cluster controler shall check all additional inputs from outside the system, the Cluster controller needs to decide whether a requested telltale is displayed for E2E miss" STYLE_REF="Requirement" ID="ID_1404407311" CREATED="1610654302938" MODIFIED="1610897196083">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Not Allocated"/>
<node TEXT="The Safety-Signal-Source shall check the additional inputs for E2E misses" STYLE_REF="Requirement" ID="ID_114212614" CREATED="1609431592950" MODIFIED="1610897196085">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety-Signal-source"/>
<node TEXT="This refers not only to the telltale request messages from the telltale requester, but also all further inputs the safety-signal source needs to decide whether the requested telltale is displayed or not, e.g. input from a HW checker element, or Image data flowing back from the display" STYLE_REF="Requirement" ID="ID_1017729133" CREATED="1610623331702" MODIFIED="1610897196086">
<attribute NAME="Type" VALUE="Information"/>
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Allocation" VALUE=""/>
</node>
</node>
</node>
<node TEXT="If the cluster controller determines an E2E miss in an additional input needed for telltale verification, the cluster controller shall transition the system into the safe state" STYLE_REF="Requirement" ID="ID_1264174165" CREATED="1610700002683" MODIFIED="1610897196087">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
<node TEXT="On E2E miss of any input to Safety-signal-source, Safety-signal-source shall request &quot;Safe state&quot; from the safety-app" STYLE_REF="Requirement" ID="ID_1488369061" CREATED="1609431616377" MODIFIED="1610897196089">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety-Signal-source"/>
</node>
</node>
</node>
<node TEXT="The Instrument cluster shall display the requested telltale or transition to the safe state" STYLE_REF="Requirement" ID="ID_1579674255" CREATED="1609154144484" MODIFIED="1610897196090">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="FSR"/>
<node TEXT="We implement this by splitting into a QM path rendering the Display and a Safety path checking whether the requested telltale is shown or not" STYLE_REF="Requirement" ID="ID_1284231708" CREATED="1609429731106" MODIFIED="1610897196092">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="The Instrument Cluster shall render the cluster display image within 70ms of the instrument Cluster receiving the message" STYLE_REF="Requirement" ID="ID_205232490" CREATED="1610125989843" MODIFIED="1610897196095">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="QM[B]"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
<node TEXT="The QT app shall render the image within 70ms of the cluster controller receiving the message" STYLE_REF="Requirement" ID="ID_499334358" CREATED="1609430956929" MODIFIED="1610897196100">
<attribute NAME="ASIL" VALUE="QM[B]"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="QT Application"/>
</node>
</node>
<node TEXT="The Instrument Cluster shall determine, whether the requested telltale is displayed" STYLE_REF="Requirement" ID="ID_874940663" CREATED="1610203034090" MODIFIED="1610897196102">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
<node TEXT="Safety-signal source part of the control flow" STYLE_REF="Requirement" ID="ID_994205752" CREATED="1609431861307" MODIFIED="1610897196103">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="The safety-signal-source shall determine, whether the requested telltale is shown" STYLE_REF="Requirement" ID="ID_745377459" CREATED="1609429843792" MODIFIED="1610897196105">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety-Signal-source"/>
</node>
<node TEXT="If the requested telltale is not shown, the Safety-signal-source shall request &quot;Safe state&quot; from the safety app." STYLE_REF="Requirement" ID="ID_1088404633" CREATED="1609431675235" MODIFIED="1610897196106">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
</node>
<node TEXT="The safety-signal source shall cyclically send the safety status message to the safety app" STYLE_REF="Requirement" ID="ID_382560048" CREATED="1609431951169" MODIFIED="1610897196107">
<attribute NAME="ASIL" VALUE="QM[B]"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety-Signal-source"/>
</node>
<node TEXT="Communication from the Safety signal source to the Safety App shall be E2E protected" STYLE_REF="Requirement" ID="ID_70275415" CREATED="1609432146314" MODIFIED="1610897196108">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety-Signal-source"/>
<node TEXT="We don&apos;t specify this in all detail here, Message counter and CRC is needed" STYLE_REF="Requirement" ID="ID_1259502493" CREATED="1610617507393" MODIFIED="1610897196109">
<attribute NAME="Type" VALUE="Information"/>
</node>
</node>
<node TEXT="The results of the Safety signal source workload shall deterministically depend on the inputs" STYLE_REF="Requirement" ID="ID_1639133793" CREATED="1609432250721" MODIFIED="1610897196109">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
<node TEXT="This implies freedom from spatial interference between the safety-signal-source / safety app and the rest of the (Operating) system, if taken at face value. The formulation is deliberate, the Architecture Workgroup is analysing all potential ways such interference could happen, we then revisit this requirement to refine it regarding safety mechanisms on the application level handling the determined interference scenarios, where possible to avoid putting undue burden on the kernel." STYLE_REF="Requirement" ID="ID_220738134" CREATED="1609432347425" MODIFIED="1610897196110">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Hardware faults are out of scope, see assumptions" STYLE_REF="Requirement" ID="ID_1937203672" CREATED="1609432645149" MODIFIED="1610897196111">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Temporal interference is not relevant here, since the watchdog transitions the system into the safe state, if execution takes too long." STYLE_REF="Requirement" ID="ID_991487171" CREATED="1609432698331" MODIFIED="1610897196112">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
</node>
</node>
<node TEXT="If the requested telltale is not displayed, the instrument cluster shall transition the system to the safe state by not triggering the external watchdog" STYLE_REF="Requirement" ID="ID_1791854442" CREATED="1610203085715" MODIFIED="1610897196112">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
<node TEXT="Safety App portion of the Control Flow" STYLE_REF="Requirement" ID="ID_538932640" CREATED="1609431908398" MODIFIED="1610897196113">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="The Safety App shall check the Communication from Safety Signal Source for E2E misses" STYLE_REF="Requirement" ID="ID_563434302" CREATED="1609432787604" MODIFIED="1610897196114">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety App"/>
</node>
<node TEXT="The Safety App shall pet the external watchdog, if and only if the cyclic message from the safety signal source passes the E2E check and does not request &quot;safe state&quot;" STYLE_REF="Requirement" ID="ID_1726916053" CREATED="1609432835072" MODIFIED="1610897196115">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
<node TEXT="Watchdog Open" STYLE_REF="KSR" ID="ID_1916374354" CREATED="1613733822060" MODIFIED="1613733928099">
<attribute_layout NAME_WIDTH="28.499999150633837 pt" VALUE_WIDTH="28.499999150633837 pt"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Watchdog Timeout Integrity" STYLE_REF="KSR" ID="ID_1446056246" CREATED="1613395692228" MODIFIED="1613481172006">
<attribute_layout NAME_WIDTH="59.999998211860714 pt" VALUE_WIDTH="59.999998211860714 pt"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
</node>
<node TEXT="Watchdog Timeout Setting" STYLE_REF="KSR" ID="ID_679959584" CREATED="1613399770698" MODIFIED="1613399847920">
<attribute_layout NAME_WIDTH="50.9999984800816 pt" VALUE_WIDTH="50.9999984800816 pt"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Watchdog Write" STYLE_REF="KSR" ID="ID_1445626914" CREATED="1613400504144" MODIFIED="1613400600951">
<attribute_layout NAME_WIDTH="28.499999150633837 pt" VALUE_WIDTH="28.499999150633837 pt"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Starting a new process" STYLE_REF="KSR" ID="ID_613007855" CREATED="1614008728724" MODIFIED="1614008767435">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Process Address Space Protection" STYLE_REF="KSR" ID="ID_1831665066" CREATED="1614008827243" MODIFIED="1614008848102">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
</node>
<node TEXT="The results of the Safety-app workload shall deterministically depend on the inputs" STYLE_REF="Requirement" ID="ID_971824356" CREATED="1609432250721" MODIFIED="1610897196116">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Safety App"/>
<node TEXT="This implies freedom from spatial interference between the safety-signal-source / safety app and the rest of the (Operating) system, if taken at face value. The formulation is deliberate, the Architecture Workgroup is analysing all potential ways such interference could happen, we then revisit this requirement to refine it regarding safety mechanisms on the application level handling the determined interference scenarios, where possible to avoid putting undue burden on the kernel." STYLE_REF="Requirement" ID="ID_1459030927" CREATED="1609432347425" MODIFIED="1610897196116">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Hardware faults are out of scope, see assumptions" STYLE_REF="Requirement" ID="ID_722885474" CREATED="1609432645149" MODIFIED="1610897196118">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Temporal interference is not relevant here, since the watchdog transitions the system into the safe state, if execution takes too long." STYLE_REF="Requirement" ID="ID_560329904" CREATED="1609432698331" MODIFIED="1610897196118">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node ID="ID_1735022659" CONTENT_ID="ID_1916374354"/>
<node ID="ID_1502095543" CONTENT_ID="ID_1446056246"/>
<node ID="ID_1760564383" CONTENT_ID="ID_679959584"/>
<node ID="ID_420297400" CONTENT_ID="ID_1445626914"/>
<node TEXT="Integrity of printf()" STYLE_REF="KSR" ID="ID_1147809717" CREATED="1614001091699" MODIFIED="1614001203367">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Integrity of mkfifo()" STYLE_REF="KSR" ID="ID_672069126" CREATED="1614001328517" MODIFIED="1614001556558">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Integrity of pipe open()" STYLE_REF="KSR" ID="ID_388336659" CREATED="1614001341827" MODIFIED="1614001556562">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Integrity of perror()" STYLE_REF="KSR" ID="ID_93094832" CREATED="1614001351500" MODIFIED="1614001556563">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="integrity of pipe read()" STYLE_REF="KSR" ID="ID_1657695425" CREATED="1614001357163" MODIFIED="1614001556564">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="integrity of system()" STYLE_REF="KSR" ID="ID_1189764432" CREATED="1614001388810" MODIFIED="1614001556565">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Integrity of usleep()" STYLE_REF="KSR" ID="ID_1535227043" CREATED="1614001401834" MODIFIED="1614001556566">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node TEXT="Integrity of fflush()" STYLE_REF="KSR" ID="ID_1314277969" CREATED="1614001411572" MODIFIED="1614001556567">
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Kernel"/>
<attribute NAME="ASIL" VALUE="B"/>
</node>
<node ID="ID_959747494" CONTENT_ID="ID_613007855"/>
<node ID="ID_1328945138" CONTENT_ID="ID_1831665066"/>
</node>
<node ID="ID_1015511320" CONTENT_ID="ID_1446056246"/>
<node ID="ID_746573891" CONTENT_ID="ID_679959584"/>
<node ID="ID_130190486" CONTENT_ID="ID_1445626914"/>
<node ID="ID_938377930" CONTENT_ID="ID_1916374354"/>
<node ID="ID_479986667" CONTENT_ID="ID_613007855"/>
<node ID="ID_192305566" CONTENT_ID="ID_1831665066"/>
</node>
<node TEXT="If the watchdog is not triggered within 200ms, it shall transition the system to the safet state" STYLE_REF="Requirement" ID="ID_998490846" CREATED="1610698701000" MODIFIED="1613400824303">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Watchdog"/>
<node TEXT="Watchdog part of the control flow" STYLE_REF="Requirement" ID="ID_1409122909" CREATED="1609431908398" MODIFIED="1610897196120">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Timing allocation considerations:&#xa;The timings for rendering and telltale verification are not safety relevant, since the watchdog transitions to the system to the safe state, if the chain takes too long." STYLE_REF="Requirement" ID="ID_1337523371" CREATED="1609430707841" MODIFIED="1610897196123">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
<node TEXT="Signal sending including rendering by QT app: 100ms. We assume the time delay between the requester sending the message, and the cluster demo receiving it is less than 30ms, leaving 70ms for the rendering" STYLE_REF="Requirement" ID="ID_865269483" CREATED="1609433185494" MODIFIED="1610897196128">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Display check inklusive WD trigger: 50ms" STYLE_REF="Requirement" ID="ID_1226012594" CREATED="1609433217685" MODIFIED="1610897196132">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node TEXT="Watchdog logic inclusive backlight killing: 50ms" STYLE_REF="Requirement" ID="ID_322365118" CREATED="1609433239777" MODIFIED="1610897196133">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
</node>
<node TEXT="The watchdog shall disable the backlight of the Cluster Display within 50ms, if it is not triggered within 150ms." STYLE_REF="Requirement" ID="ID_1266688002" CREATED="1609429267081" MODIFIED="1611331002531">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SW"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
</node>
<node ID="ID_116247226" CONTENT_ID="ID_1446056246"/>
<node ID="ID_1760172048" CONTENT_ID="ID_679959584"/>
<node ID="ID_977664596" CONTENT_ID="ID_1445626914"/>
<node ID="ID_555313202" CONTENT_ID="ID_613007855"/>
<node ID="ID_378205303" CONTENT_ID="ID_1831665066"/>
</node>
</node>
<node TEXT="The chain between Telltale request sent and display/safe state shall be less than 200ms." STYLE_REF="Requirement" ID="ID_922972509" CREATED="1609428685902" MODIFIED="1611331002531" LINK="#at(:~Sys:Telltale-requester)">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="FSR"/>
<node ID="ID_1197920546" TREE_ID="ID_1337523371">
<node ID="ID_1916288361" TREE_ID="ID_865269483"/>
<node ID="ID_980166321" TREE_ID="ID_1226012594"/>
<node ID="ID_450743490" TREE_ID="ID_322365118"/>
</node>
<node TEXT="The Telltale request message shall be sent every 200 ms" STYLE_REF="Requirement" ID="ID_190273872" CREATED="1609429199215" MODIFIED="1610897196138">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
</node>
<node ID="ID_571758931" TREE_ID="ID_205232490">
<node ID="ID_207710874" TREE_ID="ID_499334358"/>
</node>
<node TEXT="Verification of telltale shown shall be performed within 50ms" STYLE_REF="Requirement" ID="ID_1787478473" CREATED="1609430913281" MODIFIED="1610897196140">
<attribute NAME="ASIL" VALUE="QM[B]"/>
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
</node>
<node ID="ID_1679094583" TREE_ID="ID_998490846">
<node ID="ID_866134195" TREE_ID="ID_1409122909"/>
<node ID="ID_329269881" TREE_ID="ID_1337523371">
<node ID="ID_678436710" TREE_ID="ID_865269483"/>
<node ID="ID_589807630" TREE_ID="ID_1226012594"/>
<node ID="ID_27256903" TREE_ID="ID_322365118"/>
</node>
<node ID="ID_91008504" TREE_ID="ID_1266688002"/>
<node ID="ID_1627907966" CONTENT_ID="ID_1446056246"/>
<node ID="ID_285467711" CONTENT_ID="ID_679959584"/>
<node ID="ID_1843144520" CONTENT_ID="ID_1445626914"/>
<node ID="ID_1680634210" CONTENT_ID="ID_613007855"/>
<node ID="ID_268698635" CONTENT_ID="ID_1831665066"/>
</node>
</node>
</node>
<node TEXT="The system shall transition to the safe state within 100ms of the display showing an unrequested telltale for longer than 100 ms" STYLE_REF="Requirement" POSITION="right" ID="ID_888816481" CREATED="1609433460805" MODIFIED="1610897196142">
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Type" VALUE="SZ"/>
<attribute NAME="Allocation" VALUE=""/>
<node TEXT="We need to discuss this, this might not work with the frequency of 200ms we have in SZ1, it will if we relax it a little bit to around 120ms, see" STYLE_REF="Requirement" ID="ID_575915779" CREATED="1609495380329" MODIFIED="1610897196143">
<attribute NAME="ASIL" VALUE=""/>
<attribute NAME="Type" VALUE="Information"/>
</node>
<node ID="ID_1024133711" TREE_ID="ID_1780168904">
<node ID="ID_793329888" TREE_ID="ID_736988533"/>
<node ID="ID_1442215130" TREE_ID="ID_529767340"/>
<node ID="ID_9487660" TREE_ID="ID_950923064">
<node ID="ID_83965615" CONTENT_ID="ID_1340201467"/>
</node>
</node>
<node ID="ID_1726434528" TREE_ID="ID_883554261">
<node ID="ID_691030811" TREE_ID="ID_1807969240"/>
<node ID="ID_1451767216" TREE_ID="ID_199781775">
<node ID="ID_324867201" CONTENT_ID="ID_1340201467"/>
</node>
<node ID="ID_29229427" TREE_ID="ID_1213070481"/>
<node ID="ID_1988413123" TREE_ID="ID_1404407311">
<node ID="ID_571517104" TREE_ID="ID_114212614">
<node ID="ID_224260376" TREE_ID="ID_1017729133"/>
</node>
</node>
<node ID="ID_1568256292" TREE_ID="ID_1264174165">
<node ID="ID_1562874415" TREE_ID="ID_1488369061"/>
</node>
</node>
<node TEXT="The instrument cluster shall transition to the safe state within 50ms, if an unrequested telltale is displayed for more than 100 ms" STYLE_REF="Requirement" ID="ID_1967724661" CREATED="1609623643138" MODIFIED="1610897196144">
<attribute NAME="Type" VALUE="FSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<node ID="ID_464783880" TREE_ID="ID_1284231708"/>
<node ID="ID_1234093641" TREE_ID="ID_205232490">
<node ID="ID_1721681830" TREE_ID="ID_499334358"/>
</node>
<node TEXT="All Inputs the Cluster controller needs to decide whether a un requested telltale is displayed shall be E2E protected" STYLE_REF="Requirement" ID="ID_142142357" CREATED="1610654302938" MODIFIED="1610897196145">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Not Allocated"/>
</node>
<node TEXT="The Instrument Cluster shall determine, if a not requested telltale is displayed for more than 100ms" STYLE_REF="Requirement" ID="ID_1824391227" CREATED="1610203034090" MODIFIED="1610897196146">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Cluster Controller"/>
</node>
<node TEXT="If a unrequested telltale is shown for more than 100ms the instrument cluster shall transition the system to the safe state by not triggering the external watchdog" STYLE_REF="Requirement" ID="ID_1623141656" CREATED="1610660961446" MODIFIED="1610897196147">
<attribute NAME="Type" VALUE="TSR"/>
<attribute NAME="ASIL" VALUE="B"/>
<attribute NAME="Allocation" VALUE="Telltale-requester"/>
</node>
</node>
<node ID="ID_1771819379" TREE_ID="ID_922972509">
<node ID="ID_1372264395" TREE_ID="ID_1337523371">
<node ID="ID_851857056" TREE_ID="ID_865269483"/>
<node ID="ID_60352073" TREE_ID="ID_1226012594"/>
<node ID="ID_1148423018" TREE_ID="ID_322365118"/>
</node>
<node ID="ID_1374235407" TREE_ID="ID_190273872"/>
<node ID="ID_66701131" TREE_ID="ID_205232490">
<node ID="ID_139490740" TREE_ID="ID_499334358"/>
</node>
<node ID="ID_1797976261" TREE_ID="ID_1787478473"/>
<node ID="ID_666005204" TREE_ID="ID_998490846">
<node ID="ID_1120681616" TREE_ID="ID_1409122909"/>
<node ID="ID_150284297" TREE_ID="ID_1337523371">
<node ID="ID_149104298" TREE_ID="ID_865269483"/>
<node ID="ID_530521654" TREE_ID="ID_1226012594"/>
<node ID="ID_1019504250" TREE_ID="ID_322365118"/>
</node>
<node ID="ID_10230674" TREE_ID="ID_1266688002"/>
<node ID="ID_272025019" CONTENT_ID="ID_1446056246"/>
<node ID="ID_232392105" CONTENT_ID="ID_679959584"/>
<node ID="ID_435586166" CONTENT_ID="ID_1445626914"/>
<node ID="ID_1074502083" CONTENT_ID="ID_613007855"/>
<node ID="ID_16201840" CONTENT_ID="ID_1831665066"/>
</node>
</node>
</node>
<node TEXT="![system](http://www.plantuml.com/plantuml/proxy?cache=no&amp;src=https://raw.githubusercontent.com/Jochen-Kall/wg-automotive/master/AGL_cluster_demo_use_case/Architecture/Sequence-Diagram/Sequence_diagram.puml)" STYLE_REF="Requirement" POSITION="left" ID="ID_1203473398" CREATED="1609690458394" MODIFIED="1610897196150" LINK="https://github.com/Jochen-Kall/wg-automotive/blob/master/AGL_cluster_demo_use_case/Architecture/Sequence-Diagram/Sequence_diagram.md">
<attribute NAME="Type" VALUE="Information"/>
<attribute NAME="ASIL" VALUE=""/>
</node>
<node TEXT="Architecture" POSITION="left" ID="ID_1988403535" CREATED="1609881073321" MODIFIED="1610703679766" HGAP_QUANTITY="76.99999812245375 pt" VSHIFT_QUANTITY="217.49999351799508 pt">
<edge COLOR="#0000ff"/>
<node TEXT="System Architectural Elements" ID="ID_647993701" CREATED="1609881262680" MODIFIED="1610703679766">
<node TEXT="Telltale-requester" GLOBALLY_VISIBLE="true" ALIAS="Sys:Telltale-requester" ID="ID_176572829" CREATED="1609881457571" MODIFIED="1610703679766"/>
<node TEXT="Cluster Controller" ID="ID_1992686079" CREATED="1609881092272" MODIFIED="1610703679766"/>
<node TEXT="Display" ID="ID_1852933542" CREATED="1609881244846" MODIFIED="1610703679767"/>
<node TEXT="Watchdog" ID="ID_1997024973" CREATED="1609881274280" MODIFIED="1610703679767"/>
</node>
<node TEXT="SW Architectural Elements" ID="ID_1297553272" CREATED="1609881131483" MODIFIED="1610703679767">
<node TEXT="Safety-Signal-source" ID="ID_1659037005" CREATED="1609881203148" MODIFIED="1610703679767"/>
<node TEXT="Safety App" ID="ID_288712896" CREATED="1609881213584" MODIFIED="1610703679767"/>
<node TEXT="QT Application" ID="ID_1100078027" CREATED="1609881218794" MODIFIED="1610703679768"/>
<node TEXT="Kernel" ID="ID_531273097" CREATED="1610125811665" MODIFIED="1610703679768"/>
</node>
<node TEXT="HW Architectural Elements" ID="ID_983665653" CREATED="1610478772917" MODIFIED="1610703679768">
<node TEXT="Dummy" ID="ID_1264197492" CREATED="1610478968566" MODIFIED="1610703679768"/>
</node>
</node>
<node TEXT="Source code monitoring" POSITION="left" ID="ID_197895921" CREATED="1610823383222" MODIFIED="1610880288222" HGAP_QUANTITY="7.25000020116567 pt" VSHIFT_QUANTITY="35.249998949468164 pt">
<edge COLOR="#00ff00"/>
<attribute_layout NAME_WIDTH="68.99999794363981 pt" VALUE_WIDTH="248.99999257922195 pt"/>
<attribute NAME="Github link" VALUE="https://github.com/Jochen-Kall/Safety-app/"/>
<attribute NAME="revision" VALUE="8db75d886c915efc16e481e3fb63a09fd6e10eb6"/>
<node TEXT="Local Repository" ID="ID_1158912197" CREATED="1610880491236" MODIFIED="1610881409679">
<node TEXT="Safety-app" ID="ID_302688013" CREATED="1611330979915" MODIFIED="1611330979915" LINK="../../Safety-app/"/>
</node>
<node TEXT="files" ID="ID_608165453" CREATED="1610880473633" MODIFIED="1610880475535">
<node TEXT="Safety-signal-source.c" ID="ID_1396962989" CREATED="1611330996065" MODIFIED="1611330996065" LINK="../../Safety-app/Safety-signal-source.c"/>
<node TEXT="control-app.c" ID="ID_1717445163" CREATED="1611330996065" MODIFIED="1611330996065" LINK="../../Safety-app/control-app.c"/>
<node TEXT="safety-app.c" ID="ID_1170074676" CREATED="1611330996065" MODIFIED="1611330996065" LINK="../../Safety-app/safety-app.c"/>
</node>
</node>
</node>
</map>
